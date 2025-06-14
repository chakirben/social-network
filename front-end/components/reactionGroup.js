'use client';

import { useState } from 'react';
import Reaction from './reaction';
import fetchReaction from '../public/utils/fetchReaction';

export default function ReactionGroup({ itemId, itemType, userReaction, likeCount = 0, dislikeCount = 0 }) {

    const [reaction, setReaction] = useState(userReaction);
    const [likes, setLikes] = useState(likeCount);
    const [dislikes, setDislikes] = useState(dislikeCount);

    const handleReaction = async (newReaction) => {
        const current = reaction;
        const next = current === newReaction ? 0 : newReaction;

        const success = await fetchReaction(itemId, itemType, next);
        if (!success) return;

        if (current === 1) setLikes((l) => l - 1);
        if (current === -1) setDislikes((d) => d - 1);
        if (next === 1) setLikes((l) => l + 1);
        if (next === -1) setDislikes((d) => d + 1);

        setReaction(next);
    };

    return (
        <div className="reaction-group df gp12">
            <Reaction
                reactionType="like"
                count={likes}
                isReacted={reaction === 1}
                onClick={() => handleReaction(1)}
            />
            <Reaction
                reactionType="dislike"
                count={dislikes}
                isReacted={reaction === -1}
                onClick={() => handleReaction(-1)}
            />
        </div>
    );
}
