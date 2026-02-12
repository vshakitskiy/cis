import { useState } from 'react';
import { Header } from '@/components/layout/Header';
import { ProjectList } from '@/components/projects/ProjectList';
import { ProjectForm } from '@/components/projects/ProjectForm';
import { Modal } from '@/components/ui/modal';
import { Button } from '@/components/ui/button';
import { Spinner } from '@/components/ui/spinner';
import { useFetch } from '@/hooks/useFetch';
import { listProjects, createProject } from '@/api/projects';
import type { Project } from '@/types';

export function DashboardPage() {
  const { data: projects, loading, error, refetch } = useFetch<Project[]>(listProjects);
  const [showCreate, setShowCreate] = useState(false);

  const handleCreate = async (data: { name: string; description?: string; deadline?: string }) => {
    await createProject(data);
    setShowCreate(false);
    refetch();
  };

  return (
    <>
      <Header title="Dashboard" />
      <main className="flex-1 p-6">
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-semibold text-gray-900">Projects</h2>
          <Button onClick={() => setShowCreate(true)}>New Project</Button>
        </div>

        {loading && <Spinner />}
        {error && <p className="text-red-500">{error}</p>}
        {projects && <ProjectList projects={projects} />}

        <Modal open={showCreate} onClose={() => setShowCreate(false)} title="New Project">
          <ProjectForm onSubmit={handleCreate} onCancel={() => setShowCreate(false)} />
        </Modal>
      </main>
    </>
  );
}
