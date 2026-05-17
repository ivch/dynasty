/* eslint-disable */
function NavBar() {
  return (
    <nav className="guard-navbar">
      <a className="guard-brand" href="#">
        <img src="assets/logo-guard.png" alt="ЖК Династія" width="30" height="30" />
        <span>ЖК Династія</span>
      </a>
      <span className="guard-navbar__sub">Система заявок · КПП</span>
      <span className="guard-navbar__refresh">Оновлення кожні 60 с</span>
    </nav>
  );
}
window.NavBar = NavBar;
