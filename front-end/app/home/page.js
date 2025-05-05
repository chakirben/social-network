import Comment from "@/components/comments";
import "../home/comments.css"
export default async function Home() {

    let comm = await fetch("http://localhost:8080/api/GetComments?postId=2")
    let commentsData = await comm.json()
    console.log(commentsData);

    return (
        <div>
            {commentsData.map((c) => (
                <Comment key={c.id} comment={c} />
            ))}
        </div>
    )
}