"use client"
import FetchJoinToGroup from "./../../app/groups/api_and_funcs/fetch_req_join_gp"
export default function NoMyGroup({ group , onJoin }) {
     // For the button to JOIN group...
  const Jointogroup = (groupId) => {
    FetchJoinToGroup(groupId)
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
                <button id="buttontoseegroup" onClick={() => Jointogroup(group.Id)}>Join</button>
            </div>
        </div>
    );
}
