"use client"
import SideBar from "@/components/sidebar";
import GroupDetails from "./../api_and_funcs/fetch_groups_posts"
import "./../css/groups1.css"
import "./../css/creatgroup.css"
import "./../../home/home.css"
export default function DisplyGroup({ groupId , back}) {
    console.log(groupId);
    
    return (
         <div className="home">
              <SideBar />
              <div className="divallGroups">
                <GroupDetails groupId={groupId}  />
              </div>
            </div>
    )
}