import { useState } from 'react';
import type { Comment } from '@/types';
import { formatDate } from '@/lib/utils';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';

interface CommentItemProps {
  comment: Comment;
  currentUserId: number;
  onUpdate: (id: number, content: string) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
}

export function CommentItem({ comment, currentUserId, onUpdate, onDelete }: CommentItemProps) {
  const [editing, setEditing] = useState(false);
  const [content, setContent] = useState(comment.content);

  const isOwner = comment.user_id === currentUserId;

  const handleSave = async () => {
    await onUpdate(comment.id, content);
    setEditing(false);
  };

  return (
    <div className="border-b border-gray-100 py-3 last:border-0">
      <div className="flex items-center justify-between mb-1">
        <span className="text-xs text-gray-400">
          User #{comment.user_id} &middot; {formatDate(comment.created_at)}
        </span>
        {isOwner && !editing && (
          <div className="flex gap-1">
            <Button variant="ghost" className="text-xs px-2 py-1" onClick={() => setEditing(true)}>
              Edit
            </Button>
            <Button variant="ghost" className="text-xs px-2 py-1 text-red-500" onClick={() => onDelete(comment.id)}>
              Delete
            </Button>
          </div>
        )}
      </div>
      {editing ? (
        <div className="space-y-2">
          <Textarea
            className="w-full"
            rows={2}
            value={content}
            onChange={(e) => setContent(e.target.value)}
          />
          <div className="flex gap-2">
            <Button className="text-xs px-3 py-1" onClick={handleSave}>Save</Button>
            <Button variant="secondary" className="text-xs px-3 py-1" onClick={() => { setEditing(false); setContent(comment.content); }}>
              Cancel
            </Button>
          </div>
        </div>
      ) : (
        <p className="text-sm text-gray-700 whitespace-pre-wrap">{comment.content}</p>
      )}
    </div>
  );
}
