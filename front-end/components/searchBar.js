import React, { useState } from 'react';
import {FetchSearch} from '../app/home/fetch_search'
export default function SearchBar() {
    const [searchTerm, setSearchTerm] = useState('');

    const handleSearch = () => {
        console.log('Searching for:', searchTerm);
    };

    const handleKeyUp = (e) => {
        FetchSearch({searchTerm})
    };

    return (
        <div className="SearchBar">
            <div>
                <img src="./images/search.svg" alt="search icon" />
                <input
                    type="text"
                    placeholder="Type to search"
                    value={searchTerm}
                    onChange={(e) => setSearchTerm(e.target.value)}
                    onKeyUp={handleKeyUp}
                    className="searchInput"
                />
            </div>
            <button onClick={handleSearch}>Search</button>
        </div>
    );
};
