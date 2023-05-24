import React, { useEffect, useState } from "react";
import SearchBox from "./components/SearchBox";
import SearchResults from "./components/SearchResults";
import { SearchResult } from "./utils/types";
import Pagination from "./components/Pagination";
import logo from "./logo.svg";
import "./App.css";
import { backend_url } from "./utils/config";

const App: React.FC = () => {
  const [results, setResults] = useState<SearchResult[]>([]);
  const [isSearching, setIsSearching] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [resultsPerPage, setResultsPerPage] = useState(10);
  const [query, setQuery] = useState<string>("");

  const handleSearch = (searchQuery: string, page: number) => {
    setIsSearching(true);

    fetch(
      `${backend_url}/search?q=${searchQuery}&page=${page}&limit=${resultsPerPage}`
    )
      .then((response) => response.json())
      .then((data) => {
        setResults(
          data.map((item: any) => ({
            score: item.Score,
            ...item.Doc,
          }))
        );
        setIsSearching(false);
        setQuery(searchQuery);
        setCurrentPage(page);
      })
      .catch((error) => {
        console.error("Error:", error);
        setIsSearching(false);
      });
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    handleSearch(query, page);
  };

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
      {/* Add Pagination component here*/}
      <Pagination currentPage={currentPage} onPageChange={handlePageChange} />
    </div>
  );
};

export default App;
