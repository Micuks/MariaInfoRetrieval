import { renderToString } from "react-dom/server";
import "./SearchResults.css";
import React, { useEffect, useState } from "react";
import { SearchResult } from "../utils/types";
import { Prism as SyntaxHighlighter } from "react-syntax-highlighter";
import { atomDark as codeTheme } from "react-syntax-highlighter/dist/esm/styles/prism";
import { backend_url } from "../utils/config";

interface SearchResultProps {
  result: SearchResult;
}

const SearchResultItem: React.FC<SearchResultProps> = ({ result }) => {
  const [content, setContent] = useState("");
  const [contentFetched, setContentFetched] = useState<boolean>(false);
  const [isFolded, setIsFolded] = useState(true); // new state
  const toggleFold = () => {
    setIsFolded(!isFolded);
  }; // function to toggle the fold

  useEffect(() => {
    if (!isFolded && !contentFetched) {
      fetch(`${backend_url}/document?id=${result.id}`)
        .then((response) => response.json())
        .then((data) => {
          let doc = renderCodeBlocks(data.content);
          setContent(doc.body.innerHTML);
          setContentFetched(true);
        });
    }
  }, [isFolded, content, contentFetched, result]);

  const renderCodeBlocks = (content: string) => {
    let parser = new DOMParser();
    let doc = parser.parseFromString(content, "text/html");
    let codeBlocks = doc.querySelectorAll("pre");

    codeBlocks.forEach((block) => {
      let codeString = block.textContent || "";
      let highlightedCodeString = renderToString(
        <SyntaxHighlighter language="cpp" style={codeTheme}>
          {codeString}
        </SyntaxHighlighter>
      );
      block.innerHTML = highlightedCodeString;
    });

    return doc;
  };

  const handleFeedback = (resultId: string, score: number) => {
    fetch(`${backend_url}/feedback`, {
      method: "POST",
      headers: {
        "Content-TYpe": "application/json",
      },
      body: JSON.stringify({ resultId, score }),
    })
      .then((response) => response.json())
      .then((data) => console.log(data))
      .catch((error) => console.error("Error: ", error));
  };

  return (
    <div>
      <h2 onClick={toggleFold}>
        <div dangerouslySetInnerHTML={{ __html: result.title }} />
      </h2>
      {!isFolded && <div dangerouslySetInnerHTML={{ __html: content }} />}
      <p>{result.url}</p>
      <p>{result.date}</p>
      <p>Relevance: {result.score}</p>
      <button className="feedback" onClick={() => handleFeedback(result.id, 1)}>
        👍
      </button>
      <button className="feedback" onClick={() => handleFeedback(result.id, 0)}>
        👎
      </button>
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
        <SearchResultItem key={`${result.id}-${index}`} result={result} />
      ))}
    </div>
  );
};

export default SearchResults;
