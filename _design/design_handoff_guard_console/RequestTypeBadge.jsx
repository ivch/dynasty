/* eslint-disable */
function RequestTypeBadge({ rtype }) {
  const t = window.GUARD_DATA.reqTypes[rtype];
  if (!t) return null;
  return (
    <span className="guard-badge" style={{ background: t.color }}>{t.ua}</span>
  );
}
window.RequestTypeBadge = RequestTypeBadge;
