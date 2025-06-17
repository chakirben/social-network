import React from 'react';
import styles from './searchBar.module.css'; 

export default function SearchBar({ searchTerm, setSearchTerm }) {
  const handleSearch = () => {
    console.log('Searching for:', searchTerm);
  };

  const handleKeyUp = (e) => {
    if (e.key === 'Enter') {
      handleSearch();
    }
  };

  return (
    <div className={styles.SearchBar}> 
      <div>
        <img src="./images/search.svg" alt="search icon" />
        <input
          type="text"
          placeholder="Type to search"
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          onKeyUp={handleKeyUp}
          className={styles.searchInput}
        />
      </div>
    </div>
  );
}
