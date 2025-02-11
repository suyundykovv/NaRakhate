package servers

import (
	"Aitu-Bet/internal/models"
	"database/sql"
	"log"
	"math/rand"
	"time"
)

// Доступные символы на барабанах
var slotSymbols = []string{"🍒", "🍋", "🍉", "⭐", "7️⃣"}

// SpinSlots - выполняет вращение барабанов и проверяет выигрыш
func (s *Server) SpinSlots(userID int, betAmount float64) (*models.SlotSpin, error) {
	// Проверяем баланс пользователя
	var userBalance float64
	err := s.db.QueryRow("SELECT wincash FROM users WHERE id = $1", userID).Scan(&userBalance)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if userBalance < betAmount {
		return nil, sql.ErrNoRows // Недостаточно средств
	}

	// Крутим три барабана
	rand.Seed(time.Now().UnixNano())
	reel1 := slotSymbols[rand.Intn(len(slotSymbols))]
	reel2 := slotSymbols[rand.Intn(len(slotSymbols))]
	reel3 := slotSymbols[rand.Intn(len(slotSymbols))]

	// Проверяем, выиграл ли игрок
	payoutMultiplier := calculatePayout(reel1, reel2, reel3)
	payout := betAmount * payoutMultiplier
	newBalance := userBalance - betAmount + payout

	log.Printf("Обновление баланса пользователя %d: старый баланс %.2f, ставка %.2f, выплата %.2f, новый баланс %.2f",
		userID, userBalance, betAmount, payout, newBalance)

	// Обновляем баланс пользователя
	_, err = s.db.Exec("UPDATE users SET wincash = $1 WHERE id = $2", newBalance, userID)
	if err != nil {
		return nil, err
	}

	// Сохраняем результат вращения в БД
	_, err = s.db.Exec("INSERT INTO slot_spins (user_id, bet_amount, reel1, reel2, reel3, payout) VALUES ($1, $2, $3, $4, $5, $6)",
		userID, betAmount, reel1, reel2, reel3, payout)
	if err != nil {
		return nil, err
	}

	log.Printf("Сохранение спина: user_id=%d, ставка=%.2f, барабаны=[%s, %s, %s], выплата=%.2f",
		userID, betAmount, reel1, reel2, reel3, payout)

	// Возвращаем результат игроку
	spin := &models.SlotSpin{
		UserID:    userID,
		BetAmount: betAmount,
		Reel1:     reel1,
		Reel2:     reel2,
		Reel3:     reel3,
		Payout:    payout,
	}
	return spin, nil
}

// calculatePayout - рассчитывает коэффициент выигрыша
func calculatePayout(r1, r2, r3 string) float64 {
	if r1 == r2 && r2 == r3 {
		switch r1 {
		case "7️⃣":
			return 10.0 // Джекпот
		case "⭐":
			return 5.0 // Большой выигрыш
		case "🍉":
			return 3.0
		case "🍋":
			return 2.0
		case "🍒":
			return 1.5
		}
	}
	return 0.0 // Ничего не выиграл
}
