"use client"
import Creatgroup from "@/components/groups/creatGroup"
export default function DataToCreatGroup({onCreate}){
    return (
        <div className="creatgroups">
            <Creatgroup onCreate={onCreate}/>
        </div>
    )
}