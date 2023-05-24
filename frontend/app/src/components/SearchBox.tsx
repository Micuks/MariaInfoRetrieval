import "./SearchBox.css";
import { useState } from "react";
import React from "react";

interface SearchBoxProps {
  onSearch: (query: string, page: number) => void;
  onImageSearch: (image: File | null) => void;
  isSearching: boolean;
  setIsSearching: (isSearching: boolean) => void;
  query: string;
  setQuery: (query: string) => void;
}

const SearchBox: React.FC<SearchBoxProps> = ({
  onSearch,
  onImageSearch,
  isSearching,
  setIsSearching,
  query,
  setQuery,
}) => {
  const [file, setFile] = useState<File | null>(null);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    setIsSearching(true);
    if (file) {
      onImageSearch(file);
    } else {
      onSearch(query, 1); // 1 as the initial page number
    }
  };

  return (
    <div className="SearchBox">
      <form onSubmit={handleSubmit}>
        {/* Search by text */}
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Search..."
        />
        {/* Search by image */}
        <input
          type="file"
          accept="image/*"
          onChange={(e) => setFile(e.target.files?.item(0) || null)}
          id="fileUpload"
        />
        <button type="submit" disabled={isSearching}>
          Search
        </button>
      </form>
    </div>
  );
};

export default SearchBox;
