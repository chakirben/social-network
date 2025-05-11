"use client"
import SideBar from "@/components/sidebar";
import GroupDetails from "./fetch_groups_posts"
import "./../css/groups1.css"
import "./../css/creatgroup.css"
import "./../../home/home.css"
import { useState , useEffect} from "react";
export default function DisplyGroup({ groupId , back}) {
    console.log("fahd",groupId);
    const [id , setid] = useState()
    useEffect(() => {
      setid(groupId);
    }, []);
    return (
         <div className="home">
              <SideBar />
              <div className="divallGroups">
                <GroupDetails groupId={id}  />
              </div>
            </div>
    )
}