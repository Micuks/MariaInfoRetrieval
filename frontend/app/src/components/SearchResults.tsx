import "./SearchResults.css";
import React, { useState } from "react";
import { SearchResult } from "../utils/types";

interface SearchResultProps {
  result: SearchResult;
}

const SearchResultItem: React.FC<SearchResultProps> = ({ result }) => {
  const [isFolded, setIsFolded] = useState(true); // new state
  const toggleFold = () => setIsFolded(!isFolded); // function to toggle the fold
  return (
    <div>
      <h2>{result.title}</h2>
      <p
        onClick={toggleFold}
        style={{ cursor: "pointer", whiteSpace: "pre-line" }}
      >
        {isFolded ? "Show content" : result.content}
      </p>
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
