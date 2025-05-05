import React, { useState } from 'react';
export default function SearchBar() {
    const [searchTerm, setSearchTerm] = useState('');

    const handleSearch = () => {
        console.log('Searching for:', searchTerm);
    };

    return (
        <div className="SearchBar">
            <img src="./images/search.svg"></img>
            <input
                type="text"
                placeholder="Type to search"
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="searchInput"
            />
        </div>
    );
};
