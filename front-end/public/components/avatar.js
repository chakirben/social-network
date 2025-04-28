export default function AvatarUpload({ avatar, onAvatarChange }) {
  const handleChange = (e) => {
    if (e.target.files && e.target.files[0]) {
      const file = e.target.files[0];
      const objectUrl = URL.createObjectURL(file);
      onAvatarChange(objectUrl, file);
    }
  };

  return (
    <div className="avatarContainer">
      <label htmlFor="avatarUpload" className="avatarLabel">
        {avatar ? (
          <img src={avatar} alt="Avatar" className="avatarImage" />
        ) : (
          <div className="avatarPlaceholder">👤</div>
        )}
      </label>
      <input
        id="avatarUpload"
        type="file"
        accept="image/*"
        onChange={handleChange}
        style={{ display: 'none' }}
      />
    </div>
  );
}
