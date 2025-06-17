"use client"
import { use, useState } from "react"
import FetchJoinToGroup from "./../../app/groups/api_and_funcs/fetch_req_join_gp"
import FetchCancelToJoingroup from "./../../app/groups/api_and_funcs/fetch_cancel_join"
export default function NoMyGroup({ group }) {
    const [pending, setpending] = useState("")

    // For the button to JOIN group...
    const Jointogroup = (groupId) => {
        const pd = FetchJoinToGroup(groupId)
        setpending(pd)
    }
    const CancelToJoingroup = (groupId) => {
        const pd = FetchCancelToJoingroup(groupId)
        setpending("")
    }
    return (
        <div className="groupc">
            <div className="groupContent">
                <h3>{group.Title}</h3>
                <p className="discriptiong">{group.Description}</p>
                <div className="usersandpost">
                    <img src="./images/users.svg"></img>
                    <p>{group.MembersCount}Users</p>
                    <img src="./images/postgroups.svg"></img>
                    <p> {group.PostCont} Posts</p>
                </div>
            </div>
            <div className="buttonjoin">
                { pending ? (
                    <>
                       <button>Pending</button>
                       <div className="cancel" onClick={() => CancelToJoingroup(group.Id)}>cancel</div>
                    </>
                ) : (
                    <button className="tertiary" onClick={() => Jointogroup(group.Id)}>Join</button>
                )
                }
            </div>
        </div>
    );
}
