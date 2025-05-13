"use client"
import SideBar from "@/components/sidebar";
import GroupDetails from "../api_and_funcs/fetch_groups_posts"
import "./../css/groups1.css"
import "./../css/creatgroup.css"
import "./../../home/home.css"
import { use } from "react";
import { useSearchParams } from 'next/navigation';

export default function DisplyGroup({ params }) {
    const { id } = use(params);
  
    const searchParams = useSearchParams();
    const title = searchParams.get('title');

    console.log("here ---->>>>", id , title);
    
    

    return (
         <div className="home">
              <SideBar />
              <GroupDetails groupId={id} title={title} />
            </div>
    )
}