import React, { useState } from "react";
import SearchBox from "./components/SearchBox";
import SearchResults from "./components/SearchResults";
import { SearchResult } from "./utils/types";
import logo from "./logo.svg";
import "./App.css";

const App: React.FC = () => {
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);

  const backend_url = "http://47.92.133.82:9011";

  const handleSearch = (query: string) => {
    setIsSearching(true);

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
        setIsSearching(false);
      })
      .catch((error) => {
        console.error("Error:", error);
        setIsSearching(false);
      });
  };

  console.debug(`isSearching: ${isSearching}`);
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logl" />
        <p>Maria Info Retrieval System</p>
        <a
          className="App-link"
          href="https://github.com/Micuks"
          target="_blank"
          rel="noopener noreferrer"
        >
          Micuks
        </a>
      </header>
      <SearchBox
        onSearch={handleSearch}
        isSearching={isSearching}
        setIsSearching={setIsSearching}
      />
      {isSearching && <p className="Notification">Searching...</p>}
      {/* Placeholder */}
      {isSearching && <div className="Placeholder"></div>}
      <SearchResults results={results} />
    </div>
  );
};

export default App;
