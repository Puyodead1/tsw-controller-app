import Link from "next/link";

export const Nav = () => (
  <nav className="navbar bg-base-100 shadow-sm">
    <div className="flex-1">
      <ul className="menu menu-horizontal px-1">
        <li>
          <Link className="link link-hover" href="/">
            Home
          </Link>
        </li>
        <li>
          <Link className="link link-hover" href="/profile-builder">
            Profile Builder
          </Link>
        </li>
        <li>
          <Link className="link link-hover" href="/docs">
            Documentation
          </Link>
        </li>
      </ul>
    </div>
  </nav>
);
