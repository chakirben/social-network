'use client'

import { useEffect, useState } from "react"

import UserCard from "@/components/userCard";
import NoMyGroup from "./groups/notmygroups";

export default function SearchTerm(Search) {


    const [searchdata, setSearchTerm] = useState(null)

    useEffect(() => {
        const FetchSearchData = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/SearchData?query=${encodeURIComponent(Search.search)}`, {
                    credentials: "include"
                });
                if (!response.ok) {
                    console.log("error to fetch SearchData");
                    return;
                }
                const Data = await response.json();
                setSearchTerm(Data);
            } catch (error) {
                console.error("Error during search fetch:", error);
            }
        };

        if (Search.search.trim() !== "") {
            FetchSearchData();
        }
    }, [Search.search]);

    console.log(searchdata);

    if (!searchdata || (!searchdata.Notfollowed && !searchdata.UnJoinGroups)) {
        return <p>Loading or not search found...</p>;
    }
    return (
        <div className="divallGroups">
            <div className="groupsmn">
                
                {searchdata.Notfollowed && searchdata.Notfollowed.length > 0 ? (
                    searchdata.Notfollowed.map(user => (
                        <UserCard key={user.id} user={user} />
                    ))
                ) : null}
                {searchdata.UnJoinGroups && searchdata.UnJoinGroups.length > 0 ? (
                    searchdata.UnJoinGroups.map(g => (
                        <NoMyGroup key={g.Id} group={g} />
                    ))
                ) : null}
            </div>

            </div>

        
    )
}
