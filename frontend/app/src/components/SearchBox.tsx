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

  console.debug(`isSearching: ${isSearching}`);
  return (
    <div className="SearchBox">
      <input
        type="text"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search..."
      />
      <button type="submit" onClick={handleSubmit}>
        Search
      </button>
    </div>
  );
};

export default SearchBox;
