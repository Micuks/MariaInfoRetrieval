import "./SearchResults.css";
import React, { useEffect, useRef, useState } from "react";
import { SearchResult } from "../utils/types";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { dark as codeTheme } from "react-syntax-highlighter/dist/esm/styles/prism";
import { createRoot } from "react-dom/client";

interface SearchResultProps {
  result: SearchResult;
}

const SearchResultItem: React.FC<SearchResultProps> = ({ result }) => {
  const contentRef = useRef<HTMLDivElement | null>(null);
  const [isFolded, setIsFolded] = useState(true); // new state
  const toggleFold = () => setIsFolded(!isFolded); // function to toggle the fold

  useEffect(() => {
    if (contentRef.current) {
      const divElement = contentRef.current;
      let parser = new DOMParser();
      let doc = parser.parseFromString(result.content, "text/html");
      let codeBlocks = doc.querySelectorAll("pre code");

      codeBlocks.forEach((block) => {
        let replacement = document.createElement("div");
        let codeString = block.textContent || "";
        createRoot(replacement).render(
          <SyntaxHighlighter language="cpp" style={codeTheme}>
            {codeString}
          </SyntaxHighlighter>
        );
        block.parentElement?.replaceChild(replacement, block);
      });
      divElement.innerHTML = doc.body.innerHTML;
    }
  }, [result]);

  return (
    <div>
      <h2 onClick={toggleFold}>
        <div dangerouslySetInnerHTML={{ __html: result.title }} />
      </h2>
      {!isFolded && (
        // <div dangerouslySetInnerHTML={{ __html: result.content }} />
        <div ref={contentRef} />
      )}
      <p>{result.url}</p>
      <p>{result.date}</p>
      <p>Relevance: {result.score}</p>
    </div>
  );
};

interface SearchResultsProps {
  results: SearchResult[];
}

const SearchResults: React.FC<SearchResultsProps> = ({ results }) => {
  console.debug(`Generate SearchResults for`);
  console.debug(results);

  return (
    <div>
      {results.map((result, index) => (
        <SearchResultItem key={index} result={result} />
      ))}
    </div>
  );
};

export default SearchResults;
