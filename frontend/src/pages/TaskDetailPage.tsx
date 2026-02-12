import { useState, useCallback } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { Header } from '@/components/layout/Header';
import { TaskForm } from '@/components/tasks/TaskForm';
import { StatusBadge } from '@/components/tasks/StatusBadge';
import { PriorityBadge } from '@/components/tasks/PriorityBadge';
import { CommentList } from '@/components/comments/CommentList';
import { AttachmentList } from '@/components/attachments/AttachmentList';
import { TimeEntryList } from '@/components/time-entries/TimeEntryList';
import { Modal } from '@/components/ui/modal';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';
import { useFetch } from '@/hooks/useFetch';
import { useAuth } from '@/context/AuthContext';
import { getTask, updateTask, deleteTask } from '@/api/tasks';
import { listComments, createComment, updateComment, deleteComment } from '@/api/comments';
import { listAttachments, uploadAttachment, deleteAttachment } from '@/api/attachments';
import { listTimeEntries, createTimeEntry, deleteTimeEntry } from '@/api/time-entries';
import { formatDate } from '@/lib/utils';
import type { Task, Comment, Attachment, TimeEntry } from '@/types';

type Tab = 'comments' | 'attachments' | 'time';

export function TaskDetailPage() {
  const { id, taskId: taskIdParam } = useParams<{ id: string; taskId: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const projectId = Number(id);
  const taskId = Number(taskIdParam);

  const { data: task, loading, refetch: refetchTask } = useFetch<Task>(
    useCallback(() => getTask(projectId, taskId), [projectId, taskId]),
  );
  const { data: comments, refetch: refetchComments } = useFetch<Comment[]>(
    useCallback(() => listComments(projectId, taskId), [projectId, taskId]),
  );
  const { data: attachments, refetch: refetchAttachments } = useFetch<Attachment[]>(
    useCallback(() => listAttachments(projectId, taskId), [projectId, taskId]),
  );
  const { data: timeEntries, refetch: refetchTime } = useFetch<TimeEntry[]>(
    useCallback(() => listTimeEntries(projectId, taskId), [projectId, taskId]),
  );

  const [tab, setTab] = useState<Tab>('comments');
  const [showEdit, setShowEdit] = useState(false);

  const handleUpdate = async (data: Parameters<typeof updateTask>[2]) => {
    await updateTask(projectId, taskId, data);
    setShowEdit(false);
    refetchTask();
  };

  const handleDelete = async () => {
    if (!confirm('Delete this task?')) return;
    await deleteTask(projectId, taskId);
    navigate(`/projects/${projectId}`);
  };

  // Comment handlers
  const handleAddComment = async (content: string) => {
    await createComment(projectId, taskId, { content });
    refetchComments();
  };
  const handleUpdateComment = async (commentId: number, content: string) => {
    await updateComment(projectId, taskId, commentId, { content });
    refetchComments();
  };
  const handleDeleteComment = async (commentId: number) => {
    await deleteComment(projectId, taskId, commentId);
    refetchComments();
  };

  // Attachment handlers
  const handleUpload = async (file: File) => {
    await uploadAttachment(projectId, taskId, file);
    refetchAttachments();
  };
  const handleDeleteAttachment = async (attId: number) => {
    await deleteAttachment(projectId, taskId, attId);
    refetchAttachments();
  };

  // Time entry handlers
  const handleAddTime = async (data: { minutes: number; description?: string; date: string }) => {
    await createTimeEntry(projectId, taskId, data);
    refetchTime();
  };
  const handleDeleteTime = async (entryId: number) => {
    await deleteTimeEntry(projectId, taskId, entryId);
    refetchTime();
  };

  if (loading) return <><Header title="Task" /><Spinner /></>;

  const tabs: { key: Tab; label: string; count: number }[] = [
    { key: 'comments', label: 'Comments', count: comments?.length ?? 0 },
    { key: 'attachments', label: 'Attachments', count: attachments?.length ?? 0 },
    { key: 'time', label: 'Time Log', count: timeEntries?.length ?? 0 },
  ];

  return (
    <>
      <Header title={task?.title ?? 'Task'} />
      <main className="flex-1 p-6 space-y-6">
        <div className="text-sm text-gray-500">
          <Link to={`/projects/${projectId}`} className="hover:underline text-blue-600">
            Back to project
          </Link>
        </div>

        {task && (
          <div className="bg-white rounded-lg border border-gray-200 p-5 space-y-4">
            <div className="flex items-start justify-between">
              <div>
                <h2 className="text-xl font-semibold text-gray-900">{task.title}</h2>
                {task.description && (
                  <p className="text-gray-600 mt-1">{task.description}</p>
                )}
              </div>
              <div className="flex gap-2">
                <Button variant="ghost" onClick={() => setShowEdit(true)}>Edit</Button>
                <Button variant="destructive" onClick={handleDelete}>Delete</Button>
              </div>
            </div>
            <div className="flex items-center gap-3 flex-wrap">
              <StatusBadge status={task.status} />
              <PriorityBadge priority={task.priority} />
              {task.deadline && (
                <span className="text-sm text-gray-500">Due {formatDate(task.deadline)}</span>
              )}
              {task.assignee_id && (
                <span className="text-sm text-gray-500">Assigned to user #{task.assignee_id}</span>
              )}
            </div>
          </div>
        )}

        {/* Tabs */}
        <div className="border-b border-gray-200">
          <div className="flex gap-4">
            {tabs.map((t) => (
              <button
                key={t.key}
                className={`pb-2 text-sm font-medium border-b-2 transition-colors cursor-pointer ${
                  tab === t.key
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
                onClick={() => setTab(t.key)}
              >
                {t.label} ({t.count})
              </button>
            ))}
          </div>
        </div>

        {tab === 'comments' && comments && (
          <CommentList
            comments={comments}
            currentUserId={user?.user_id ?? 0}
            onAdd={handleAddComment}
            onUpdate={handleUpdateComment}
            onDelete={handleDeleteComment}
          />
        )}

        {tab === 'attachments' && attachments && (
          <AttachmentList
            attachments={attachments}
            projectId={projectId}
            taskId={taskId}
            onUpload={handleUpload}
            onDelete={handleDeleteAttachment}
          />
        )}

        {tab === 'time' && timeEntries && (
          <TimeEntryList
            entries={timeEntries}
            onAdd={handleAddTime}
            onDelete={handleDeleteTime}
          />
        )}

        <Modal open={showEdit} onClose={() => setShowEdit(false)} title="Edit Task">
          {task && (
            <TaskForm
              initial={task}
              onSubmit={handleUpdate}
              onCancel={() => setShowEdit(false)}
            />
          )}
        </Modal>
      </main>
    </>
  );
}
