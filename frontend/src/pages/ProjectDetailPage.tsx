import { useState, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Header } from '@/components/layout/Header';
import { TaskBoard } from '@/components/tasks/TaskBoard';
import { TaskList } from '@/components/tasks/TaskList';
import { TaskForm } from '@/components/tasks/TaskForm';
import { ProjectForm } from '@/components/projects/ProjectForm';
import { ProjectReport } from '@/components/reports/ProjectReport';
import { Modal } from '@/components/ui/modal';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';
import { useFetch } from '@/hooks/useFetch';
import { getProject, updateProject, deleteProject } from '@/api/projects';
import { listTasks, createTask, updateTask } from '@/api/tasks';
import { getProjectReport } from '@/api/reports';
import type { Project, Task, TaskStatus, ProjectReport as ProjectReportType } from '@/types';

export function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const projectId = Number(id);

  const { data: project, loading: projLoading, refetch: refetchProject } = useFetch<Project>(
    useCallback(() => getProject(projectId), [projectId]),
  );
  const { data: tasks, loading: tasksLoading, refetch: refetchTasks } = useFetch<Task[]>(
    useCallback(() => listTasks(projectId), [projectId]),
  );
  const { data: report, refetch: refetchReport } = useFetch<ProjectReportType>(
    useCallback(() => getProjectReport(projectId), [projectId]),
  );

  const [view, setView] = useState<'board' | 'list'>('board');
  const [showCreateTask, setShowCreateTask] = useState(false);
  const [showEditProject, setShowEditProject] = useState(false);

  const handleCreateTask = async (data: { title: string; description?: string; priority?: string; deadline?: string }) => {
    await createTask(projectId, data as Parameters<typeof createTask>[1]);
    setShowCreateTask(false);
    refetchTasks();
    refetchReport();
  };

  const handleUpdateProject = async (data: { name: string; description?: string; deadline?: string }) => {
    await updateProject(projectId, data);
    setShowEditProject(false);
    refetchProject();
  };

  const handleDeleteProject = async () => {
    if (!confirm('Delete this project? This cannot be undone.')) return;
    await deleteProject(projectId);
    navigate('/');
  };

  const handleStatusChange = async (taskId: number, newStatus: TaskStatus) => {
    const task = tasks?.find((t) => t.id === taskId);
    if (!task) return;
    await updateTask(projectId, taskId, {
      title: task.title,
      description: task.description ?? undefined,
      assignee_id: task.assignee_id,
      status: newStatus,
      priority: task.priority,
      deadline: task.deadline ?? undefined,
    });
    refetchTasks();
    refetchReport();
  };

  if (projLoading || tasksLoading) return <><Header title="Project" /><Spinner /></>;

  return (
    <>
      <Header title={project?.name ?? 'Project'} />
      <main className="flex-1 p-6 space-y-6">
        {report && <ProjectReport report={report} />}

        <div className="flex items-center justify-between">
          <div className="flex gap-2">
            <Button
              variant={view === 'board' ? 'default' : 'secondary'}
              onClick={() => setView('board')}
            >
              Board
            </Button>
            <Button
              variant={view === 'list' ? 'default' : 'secondary'}
              onClick={() => setView('list')}
            >
              List
            </Button>
          </div>
          <div className="flex gap-2">
            <Button variant="ghost" onClick={() => setShowEditProject(true)}>Edit Project</Button>
            <Button variant="destructive" onClick={handleDeleteProject}>Delete</Button>
            <Button onClick={() => setShowCreateTask(true)}>New Task</Button>
          </div>
        </div>

        {tasks && view === 'board' && (
          <TaskBoard tasks={tasks} projectId={projectId} onStatusChange={handleStatusChange} />
        )}
        {tasks && view === 'list' && <TaskList tasks={tasks} projectId={projectId} />}

        <Modal open={showCreateTask} onClose={() => setShowCreateTask(false)} title="New Task">
          <TaskForm onSubmit={handleCreateTask} onCancel={() => setShowCreateTask(false)} />
        </Modal>

        <Modal open={showEditProject} onClose={() => setShowEditProject(false)} title="Edit Project">
          {project && (
            <ProjectForm
              initial={project}
              onSubmit={handleUpdateProject}
              onCancel={() => setShowEditProject(false)}
            />
          )}
        </Modal>
      </main>
    </>
  );
}
