'use client'

import { useEffect, useState } from "react"
import { UserCard } from "@/components/userCard"
import { NoMyGroup } from "@/components/groups/notmygroups"
export default function SearchTerm(SearchTerm) {

    const [searchdata, setSearchTerm] = useState(null)

    useEffect(() => {
        const FetchSearchData = async () => {
            try {
                const response = await fetch(`http://localhost:8080/api/SearchData?query=${encodeURIComponent(SearchTerm)}`, {
                    credentials: "include"
                });
                if (!response.ok) {
                    console.log("error to fetch SearchData")
                    return
                }
                const Data = await response.json()
                setSearchTerm(Data)
            } catch (error) {

            }
        }
        FetchSearchData()
    }, [])

    console.log(searchdata);
    

    return (
        <>

          

        </>
    )
}

/**?
 * 
 *   <div className="usersList">
                {searchdata.NotfollowedUser.map(user => (
                    <UserCard key={user.id} user={user} />
                ))}
            </div>



            <div className="groupsmn">
                {searchdata.notMyGroups.map((group) => (
                    <NoMyGroup key={group.Id} group={group} onJoin={onJoin} />
                ))}
            </div> 



 */