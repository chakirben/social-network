import { useRouter } from 'next/navigation';

export default function Post ({pst}) {
    const router = useRouter();
    return (
        <div className="post" onClick={() => router.push(`/post/${pst.id}`)}>
            <div className="postContent">
                <div className="userData">
                    {pst.image ? <img className="pic sm" src={`http://localhost:8080/${pst.image}`} ></img>:""}
                    <h4>{pst.creator}</h4>
                    <h5>{pst.created_at}</h5>
                </div>
                <div className="content">{pst.content}</div>
                <img className="pic nrml" src={`http://localhost:8080/${pst.image}`}></img>
            </div>
        </div>
    )
}