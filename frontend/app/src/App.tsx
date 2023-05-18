import React, { useState } from "react";
import SearchBox from "./components/SearchBox";
import SearchResults from "./components/SearchResults";
import { SearchResult } from "./utils/types";
import logo from "./logo.svg";
import "./App.css";

const App: React.FC = () => {
  const [results, setResults] = useState<SearchResult[]>([]);

  const handleSearch = (query: string) => {
    fetch(`http://localhost:9011/search?q=${query}`)
      .then((response) => response.json())
      .then((data) => {
        console.debug(`Get results ${data}`);
        console.debug(data);
        setResults(
          data.map((item: any) => ({
            score: item.score,
            ...item.document,
          }))
        );
      });
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
