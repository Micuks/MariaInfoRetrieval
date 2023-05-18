import React, { useState } from "react";
import SearchBox from "./components/SearchBox";
import SearchResults from "./components/SearchResults";
import { SearchResult } from "./utils/types";
import logo from "./logo.svg";
import "./App.css";
import { response } from "express";

const App: React.FC = () => {
  const [results, setResults] = useState<SearchResult[]>([]);

  const handleSearch = (query: string) => {
    fetch(`http://localhost:9011/search?q=${query}`)
      .then((response) => response.json())
      .then((data) =>
        setResults(
          data.hits.hits.map((hit: any) => ({
            score: hit._score,
            ...hit._source,
          }))
        )
      );
  };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logl" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <SearchBox onSearch={handleSearch} />
        <SearchResults results={results} />
        <a
          className="App-link"
          href="https://github.com/Micuks"
          target="_blank"
          rel="noopener noreferrer"
        >
          Micuks
        </a>
      </header>
    </div>
  );
};

export default App;
