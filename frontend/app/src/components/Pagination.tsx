import React from "react";

interface PaginationProps {
  currentPage: number;
  onPageChange: (page: number) => void;
}

const Pagination: React.FC<PaginationProps> = ({
  currentPage,
  onPageChange,
}) => {
  return (
    <div>
      <button onClick={() => onPageChange(currentPage - 1)}>Previous</button>
      <span>{currentPage}</span>
      <button onClick={() => onPageChange(currentPage + 1)}>Next</button>
    </div>
  );
};

export default Pagination;
