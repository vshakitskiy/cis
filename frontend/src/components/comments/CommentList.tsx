import type { Comment } from '@/types';
import { CommentItem } from './CommentItem';
import { CommentForm } from './CommentForm';

interface CommentListProps {
  comments: Comment[];
  currentUserId: number;
  onAdd: (content: string) => Promise<void>;
  onUpdate: (id: number, content: string) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
}

export function CommentList({ comments, currentUserId, onAdd, onUpdate, onDelete }: CommentListProps) {
  return (
    <div className="space-y-4">
      <CommentForm onSubmit={onAdd} />
      {comments.length === 0 ? (
        <p className="text-sm text-gray-400 text-center py-4">No comments yet</p>
      ) : (
        <div>
          {comments.map((c) => (
            <CommentItem
              key={c.id}
              comment={c}
              currentUserId={currentUserId}
              onUpdate={onUpdate}
              onDelete={onDelete}
            />
          ))}
        </div>
      )}
    </div>
  );
}
