import React from 'react';

export default function SearchBar({ searchTerm, setSearchTerm }) {
    return (
        <div className="SearchBar">
            <div>
                <img src="./images/search.svg" alt="search icon" />
                <input
                    type="text"
                    placeholder="Type to search"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    className="searchInput"
                />
            </div>
        </div>
    );
}
