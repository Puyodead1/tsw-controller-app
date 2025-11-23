import { Nav } from "@/_components/Nav";
import { pages } from "@/_constants/docs.pages";
import Link from "next/link";

export default function DocsIndexPage() {
  return (
    <>
      <Nav />
      <div className="prose max-w-4xl mx-auto px-8 my-8">
        <ul>
          {Object.entries(pages).map(([slug, def]) => (
            <li key={slug}>
              <Link className="link link-hover" href={`/docs/${slug}`}>
                {def.title}
              </Link>
            </li>
          ))}
        </ul>
      </div>
    </>
  );
}
