
export default async function fetchReaction(itemId, itemType, reactionValue) {
    console.log(itemId , itemType ,  reactionValue);
    
    try {
        const response = await fetch(`/api/reaction?itemType=${itemType}&itemId=${itemId}&reactionType=${reactionValue}`, {
            method: 'POST',
            credentials: "include"
        });

        if (!response.ok) {
            console.error("Server error:", response.statusText);
            return false;
        }
        return true;
    } catch (error) {
        console.error("Network error:", error);
        return false;
    }
}
