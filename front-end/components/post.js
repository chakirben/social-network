import { useRouter } from 'next/navigation';
// import Reaction from './reaction';
import ReactionGroup from './reactionGroup';
import { timePassed } from '@/public/utils/timePassed';

export default function Post({ pst }) {
    const router = useRouter();
    return (
        <div className="post" onClick={() => router.push(`/post/${pst.id}`)}>
            <div className="content df cl gp12">
                <div className="userData">
                    {pst?.avatar ? <img className="pic sm" src={`http://localhost:8080/${pst?.avatar}`} ></img> : ""}
                    <h4>{pst?.creator}</h4>
                    <h5>{"• " +timePassed(pst?.created_at)}</h5>
                </div>
                <div className="content">{pst?.content}</div>
                {pst.image && <img className="pic nrml" src={`http://localhost:8080/${pst?.image}`}></img>}
            </div>
            {/* <div className='reactionsContainer df gp12'>
                <Reaction count={pst.like_count} itemType="post" itemId={pst.id} userReaction={pst.user_reaction} reactionType={"like"} />
                <Reaction count={pst.dislike_count} itemType="post" itemId={pst.id} userReaction={pst.user_reaction} reactionType={"dislike"} />
            </div> */}
            <ReactionGroup className="reactionsContainer df gp12" likeCount={pst.like_count} dislikeCount={pst.dislike_count} itemType="post" itemId={pst.id} userReaction={pst.user_reaction} ></ReactionGroup>

        </div>
    )
}