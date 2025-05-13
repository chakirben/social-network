"use client"
import { useState } from "react"
export default function DataToCreatGroup({onCreate , onSkip}){
        const [title, setTitle] = useState("");
        const [description, setDescription] = useState("");
    
        const handleCreate = () => {
            if (title.trim() === "" || description.trim() === "") {
                alert("Full all of the title and description..");
                return;
            }
    
            onCreate({ title, description });
        };
        return (
            <div className="creatgroups">
                <div className="inputToCreatGourp">
                    <input className="titelinp" type="text" placeholder="Title of the group..." onChange={(e) => setTitle(e.target.value)} />
                    <input className="discinp" type="text" placeholder="Description of the group..." onChange={(e) => setDescription(e.target.value)} />
                </div>
                <div className="divcreatbtn">
                    <button className="creatbtn" onClick={handleCreate} >+ Creat</button>
                    <button className="skipbtn"  onClick={onSkip}> <img src="./images/skip.svg" /> Skip</button>
                </div>
            </div>
        )
}