/* eslint-disable */
function Pager({ currentPage, totalPages, onPage }) {
  if (totalPages <= 1) return null;
  const pages = Array.from({ length: totalPages }, (_, i) => i + 1);
  return (
    <nav className="guard-pager" aria-label="Сторінки">
      <ul className="guard-pager__list">
        {pages.map((p) => (
          <li key={p} className={"guard-pager__item" + (p === currentPage ? " is-active" : "")}>
            <a
              className="guard-pager__link"
              href="#"
              onClick={(e) => { e.preventDefault(); onPage(p); }}
            >{p}</a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
window.Pager = Pager;
