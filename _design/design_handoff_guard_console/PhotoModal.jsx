/* eslint-disable */
function PhotoModal({ open, src, onClose }) {
  if (!open) return null;
  return (
    <div className="guard-modal" onClick={onClose} role="dialog" aria-modal="true">
      <div className="guard-modal__dialog" onClick={(e) => e.stopPropagation()}>
        <button className="guard-modal__close" onClick={onClose} aria-label="Закрити">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round">
            <line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/>
          </svg>
        </button>
        <div className="guard-modal__photo">
          {src && src.img ? (
            <img src={src.img} alt="Фото заявки" className="guard-modal__img" />
          ) : (
            <div className="guard-modal__placeholder">
              <svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1" strokeLinecap="round" strokeLinejoin="round" style={{opacity:0.3}}>
                <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/>
                <polyline points="21 15 16 10 5 21"/>
              </svg>
              <div style={{marginTop:12,fontSize:13}}>Фото відсутнє</div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
window.PhotoModal = PhotoModal;
