import React, { useState } from "react";
import SearchBox from "./components/SearchBox";
import SearchResults from "./components/SearchResults";
import { SearchResult } from "./utils/types";
import logo from "./logo.svg";
import "./App.css";

const App: React.FC = () => {
  const [results, setResults] = useState<SearchResult[]>([]);
  const backend_url = "http://10.128.170.37:9011";

  const handleSearch = (query: string) => {
    fetch(`${backend_url}/search?q=${query}`)
      .then((response) => response.json())
      .then((data) => {
        console.debug(data);
        setResults(
          data.map((item: any) => ({
            score: item.Score,
            ...item.Doc,
          }))
        );
      });
  };

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logl" />
        <p>Maria Info Retrieval System</p>
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
