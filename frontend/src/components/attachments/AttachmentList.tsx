import type { Attachment } from '@/types';
import { formatDate, formatBytes, downloadBlob } from '@/lib/utils';
import { downloadAttachment } from '@/api/attachments';
import { AttachmentUpload } from './AttachmentUpload';
import { Button } from '@/components/ui/button';

interface AttachmentListProps {
  attachments: Attachment[];
  projectId: number;
  taskId: number;
  onUpload: (file: File) => Promise<void>;
  onDelete: (id: number) => Promise<void>;
}

export function AttachmentList({ attachments, projectId, taskId, onUpload, onDelete }: AttachmentListProps) {
  const handleDownload = async (att: Attachment) => {
    const blob = await downloadAttachment(projectId, taskId, att.id);
    downloadBlob(blob, att.filename);
  };

  return (
    <div className="space-y-4">
      <AttachmentUpload onUpload={onUpload} />
      {attachments.length === 0 ? (
        <p className="text-sm text-gray-400 text-center py-4">No attachments yet</p>
      ) : (
        <div className="space-y-2">
          {attachments.map((att) => (
            <div key={att.id} className="flex items-center justify-between bg-gray-50 rounded-md px-3 py-2">
              <div>
                <button
                  className="text-sm text-blue-600 hover:underline font-medium cursor-pointer"
                  onClick={() => handleDownload(att)}
                >
                  {att.filename}
                </button>
                <p className="text-xs text-gray-400">
                  {formatBytes(att.size)} &middot; {formatDate(att.created_at)}
                </p>
              </div>
              <Button variant="ghost" className="text-xs text-red-500" onClick={() => onDelete(att.id)}>
                Delete
              </Button>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
