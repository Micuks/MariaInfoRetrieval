import React from "react";
import { SearchResult } from "../utils/types";
interface SearchResultsProps {
  results: SearchResult[];
}

const SearchResults: React.FC<SearchResultsProps> = ({ results }) => (
  <div>
    {results.map((result, index) => (
      <div key={index}>
        <h2>{result.title}</h2>
        <p>{result.content}</p>
        <p>{result.url}</p>
        <p>{result.date}</p>
        <p>Relecanve: {result.score}</p>
      </div>
    ))}
  </div>
);

export default SearchResults;
