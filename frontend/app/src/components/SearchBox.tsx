import "./SearchBox.css";
import { useState } from "react";
import React from "react";

interface SearchBoxProps {
  onSearch: (query: string, page: number) => void;
  isSearching: boolean;
  setIsSearching: (isSearching: boolean) => void;
}

const SearchBox: React.FC<SearchBoxProps> = ({
  onSearch,
  isSearching,
  setIsSearching,
}) => {
  const [query, setQuery] = useState<string>("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsSearching(true);
    onSearch(query, 1); // 1 as the initial page number
  };

  const handleKeyDown = (key: string) => {
    if (key === "Enter") {
      setIsSearching(true);
      onSearch(query, 1);
    }
  };

  console.debug(`isSearching: ${isSearching}`);
  return (
    <div className="SearchBox">
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        onKeyDown={(e) => handleKeyDown(e.key)}
        placeholder="Search..."
      />
      <button type="submit" onClick={handleSubmit} disabled={isSearching}>
        Search
      </button>
    </div>
  );
};

export default SearchBox;
