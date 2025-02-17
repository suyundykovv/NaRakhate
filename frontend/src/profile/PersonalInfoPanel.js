import React from "react";
import { Camera, X } from "lucide-react";

export default function PersonalInfoPanel({ show, onClose, profileData, setProfileData }) {
  const handleAvatarChange = (e) => {
    const file = e.target.files?.[0];
    if (file) {
      const imageUrl = URL.createObjectURL(file);
      setProfileData(prev => ({
        ...prev,
        avatarUrl: imageUrl
      }));
    }
  };

  const handlePersonalInfoSubmit = (e) => {
    e.preventDefault();
    onClose();
  };

  if (!show) return null;

  return (
    <div className="info-panel-overlay">
      <div className="info-panel">
        <div className="info-panel-header">
          <h3>Personal Information</h3>
          <button 
            className="close-button"
            onClick={onClose}
          >
            <X size={24} />
          </button>
        </div>
        <div className="info-panel-content">
          <form onSubmit={handlePersonalInfoSubmit}>
            <div className="avatar-upload">
              <div className="avatar-preview">
                <img
                  src={profileData.avatarUrl || "https://via.placeholder.com/100x100"}
                  alt="Avatar preview"
                />
                <label className="avatar-upload-button">
                  <Camera size={20} />
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleAvatarChange}
                    className="hidden"
                  />
                </label>
              </div>
            </div>

            <div className="form-group">
              <label htmlFor="firstName">First Name</label>
              <input
                type="text"
                id="firstName"
                value={profileData.firstName}
                onChange={(e) => setProfileData(prev => ({
                  ...prev,
                  firstName: e.target.value
                }))}
                className="form-input"
                placeholder="Enter your first name"
              />
            </div>

            <div className="form-group">
              <label htmlFor="lastName">Last Name</label>
              <input
                type="text"
                id="lastName"
                value={profileData.lastName}
                onChange={(e) => setProfileData(prev => ({
                  ...prev,
                  lastName: e.target.value
                }))}
                className="form-input"
                placeholder="Enter your last name"
              />
            </div>

            <button type="submit" className="submit-button">
              Save Changes
            </button>
          </form>
        </div>
      </div>
    </div>
  );
}