import React from "react";
import { SearchResult } from "../utils/types";
interface SearchResultsProps {
  results: SearchResult[];
}

const SearchResults: React.FC<SearchResultsProps> = ({ results }) => {
  console.debug(`Generate SearchResults for ${results}`);
  return (
    <div>
      {results.map((result, index) => (
        <div key={index}>
          <h2>{result.title}</h2>
          <p>{result.content}</p>
          <p>{result.url}</p>
          <p>{result.date}</p>
          <p>Relevance: {result.score}</p>
        </div>
      ))}
    </div>
  );
};

export default SearchResults;
