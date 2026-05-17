/* eslint-disable */
function RequestRow({ req, onAction, onImageClick }) {
  const date = new Date(req.time * 1000);
  const dateStr = date.toLocaleString("uk-UA", {
    day: "2-digit", month: "2-digit", year: "numeric",
    hour: "2-digit", minute: "2-digit",
  });

  const isClosed   = req.status === "closed";
  const actionLabel = isClosed ? "Відкрити" : "Закрити";
  const actionClass = isClosed ? "guard-btn--secondary" : "guard-btn--success";
  const nextStatus  = isClosed ? "new" : "closed";

  return (
    <tr className={isClosed ? "is-closed" : ""}>
      <td className="guard-row__main">
        <div className="guard-row__head">
          <RequestTypeBadge rtype={req.rtype} />
          <span className="guard-row__when">{dateStr}</span>
          {isClosed && <span className="guard-row__status-closed">Закрито</span>}
        </div>
        <div className="guard-row__addr">
          <strong>{req.address} <span className="guard-row__apt">#{req.apartment}</span></strong>
          <span className="guard-row__dim">{req.user_name} · {req.phone}</span>
        </div>
        {req.description && (
          <div className="guard-row__desc">{req.description}</div>
        )}
        {req.images && req.images.length > 0 && (
          <div className="guard-row__imgs">
            {req.images.map((img, i) => (
              <button
                key={i}
                type="button"
                className="guard-thumb"
                onClick={() => onImageClick(img)}
                aria-label="Переглянути фото"
              >
                <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                  <rect x="3" y="3" width="18" height="18" rx="2"/><circle cx="8.5" cy="8.5" r="1.5"/>
                  <polyline points="21 15 16 10 5 21"/>
                </svg>
              </button>
            ))}
          </div>
        )}
      </td>
      <td className="guard-row__actions">
        <button
          type="button"
          className={"guard-btn " + actionClass}
          onClick={() => onAction(req.id, nextStatus)}
        >{actionLabel}</button>
      </td>
    </tr>
  );
}
window.RequestRow = RequestRow;
