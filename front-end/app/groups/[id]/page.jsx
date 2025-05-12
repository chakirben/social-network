"use client"
import SideBar from "@/components/sidebar";
import GroupDetails from "./fetch_groups_posts"
import "./../css/groups1.css"
import "./../css/creatgroup.css"
import "./../../home/home.css"
import { use } from "react";
export default function DisplyGroup({ params }) {
    const { id } = use(params); 
    console.log("fahd jj", id);
    return (
         <div className="home">
              <SideBar />
              <div className="divallGroups">
                <GroupDetails groupId={id}/>
              </div>
            </div>
    )
}