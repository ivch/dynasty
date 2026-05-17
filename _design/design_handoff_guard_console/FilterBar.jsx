/* eslint-disable */
function FilterBar({ statusFilter, typeFilter, aptQuery, onToggle, onSearch, onClearSearch }) {
  return (
    <div className="guard-filterbar">
      <div className="guard-pills">
        <button
          className={"guard-pill" + (typeFilter === "kpp" ? " is-active" : "")}
          onClick={() => onToggle("type")}
        >Тільки для КПП</button>
        <button
          className={"guard-pill" + (statusFilter === "new" ? " is-active" : "")}
          onClick={() => onToggle("status")}
        >Тільки відкриті</button>
      </div>
      <div className="guard-search">
        <div className="guard-input-group">
          <span className="guard-input-icon">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><rect x="3" y="3" width="7" height="11" rx="1"/><path d="M10 7h11M10 11h7"/></svg>
          </span>
          <input
            type="text"
            className="guard-input"
            placeholder="квартира"
            value={aptQuery}
            maxLength={4}
            onChange={(e) => onSearch(e.target.value)}
          />
          {aptQuery && (
            <button className="guard-input-clear" onClick={onClearSearch} title="Очистити">✕</button>
          )}
        </div>
      </div>
    </div>
  );
}
window.FilterBar = FilterBar;
