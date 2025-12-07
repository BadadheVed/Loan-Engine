"use client";

import { useState, useRef } from "react";

type HealthStatus = "idle" | "success" | "error" | "loading";

function HealthCheckButton() {
  const [status, setStatus] = useState<HealthStatus>("idle");

  const checkHealth = async () => {
    setStatus("loading");
    try {
      const res = await fetch(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/health`
      );
      if (res.ok) {
        setStatus("success");
      } else {
        setStatus("error");
      }
    } catch {
      setStatus("error");
    }
  };

  let color = "bg-gray-200 text-gray-600";
  let icon = null;
  if (status === "success") {
    color = "bg-green-500 text-white";
    icon = <span className="ml-2">✔️</span>;
  } else if (status === "error") {
    color = "bg-red-500 text-white";
    icon = <span className="ml-2">❌</span>;
  } else if (status === "loading") {
    color = "bg-blue-500 text-white";
    icon = <span className="ml-2 animate-spin">⏳</span>;
  }

  return (
    <button
      onClick={checkHealth}
      className={`flex items-center px-4 py-2 rounded-xl font-medium transition-all duration-200 ${color}`}
      disabled={status === "loading"}
      title="Check backend health"
    >
      Health Check
      {icon}
    </button>
  );
}
import Image from "next/image";

export default function Dashboard() {
  const [file, setFile] = useState<File | null>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [summary, setSummary] = useState<{
    records_added: number;
    records_failed: number;
    records_skipped: number;
    duplicate_email_count: number;
  } | null>(null);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      setFile(e.dataTransfer.files[0]);
      setError(null);
      setSummary(null);
    }
  };

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
      setError(null);
      setSummary(null);
    }
  };

  const uploadFile = async () => {
    if (!file) return;

    setIsLoading(true);
    setError(null);

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_BACKEND_URL}/api/uploadcsv`,
        {
          method: "POST",
          body: formData,
        }
      );

      if (!response.ok) {
        throw new Error("Failed to upload file");
      }

      const data = await response.json();
      setSummary(data);
    } catch (err) {
      setError("An error occurred while uploading. Please try again.");
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-black text-gray-900 dark:text-gray-100 p-8 font-sans">
      <div className="max-w-4xl mx-auto space-y-8">
        <header>
          <h1 className="text-3xl font-bold tracking-tight">
            CSV Upload Dashboard
          </h1>
          <p className="text-gray-500 dark:text-gray-400 mt-2">
            Upload your user data CSV to process and analyze.
          </p>
        </header>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
          {/* Upload Section */}
          <section className="space-y-4">
            <div
              className={`border-2 border-dashed rounded-2xl p-8 transition-all duration-200 ease-in-out flex flex-col items-center justify-center text-center cursor-pointer ${
                isDragging
                  ? "border-blue-500 bg-blue-50 dark:bg-blue-900/20"
                  : "border-gray-300 dark:border-gray-700 hover:border-gray-400 dark:hover:border-gray-600"
              }`}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={() => fileInputRef.current?.click()}
            >
              <input
                type="file"
                ref={fileInputRef}
                className="hidden"
                accept=".csv"
                onChange={handleFileSelect}
              />
              <div className="w-12 h-12 mb-4 text-gray-400">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m3.75 9v6m3-3H9m1.5-12H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z"
                  />
                </svg>
              </div>
              <p className="text-sm font-medium">
                {file ? file.name : "Drag & drop or click to select CSV"}
              </p>
              {file && (
                <p className="text-xs text-gray-400 mt-1">
                  {(file.size / 1024).toFixed(1)} KB
                </p>
              )}
            </div>

            <button
              onClick={uploadFile}
              disabled={!file || isLoading}
              className={`w-full py-3 px-4 rounded-xl font-medium transition-all duration-200 ${
                !file || isLoading
                  ? "bg-gray-200 dark:bg-gray-800 text-gray-400 cursor-not-allowed"
                  : "bg-black dark:bg-white text-white dark:text-black hover:opacity-90 active:scale-[0.98]"
              }`}
            >
              {isLoading ? "Processing..." : "Upload CSV"}
            </button>

            {error && (
              <div className="p-4 rounded-xl bg-red-50 dark:bg-red-900/20 text-red-600 dark:text-red-400 text-sm">
                {error}
              </div>
            )}
          </section>

          {/* Summary Section + Health Check */}
          <section className="space-y-4">
            <div className="flex items-center justify-between">
              <h2 className="text-xl font-semibold">Processing Summary</h2>
              <HealthCheckButton />
            </div>
            {summary ? (
              <div className="grid grid-cols-2 gap-4">
                <SummaryCard
                  label="Added"
                  value={summary.records_added}
                  color="text-green-600 dark:text-green-400"
                  bg="bg-green-50 dark:bg-green-900/20"
                />
                <SummaryCard
                  label="Failed"
                  value={summary.records_failed}
                  color="text-red-600 dark:text-red-400"
                  bg="bg-red-50 dark:bg-red-900/20"
                />
                <SummaryCard
                  label="Skipped"
                  value={summary.records_skipped}
                  color="text-yellow-600 dark:text-yellow-400"
                  bg="bg-yellow-50 dark:bg-yellow-900/20"
                />
                <SummaryCard
                  label="Duplicates"
                  value={summary.duplicate_email_count}
                  color="text-blue-600 dark:text-blue-400"
                  bg="bg-blue-50 dark:bg-blue-900/20"
                />
              </div>
            ) : (
              <div className="h-48 rounded-2xl bg-gray-100 dark:bg-zinc-900 flex items-center justify-center text-gray-400 text-sm">
                Awaiting upload...
              </div>
            )}
          </section>
        </div>
      </div>
    </div>
  );
}

function SummaryCard({
  label,
  value,
  color,
  bg,
}: {
  label: string;
  value: number;
  color: string;
  bg: string;
}) {
  return (
    <div
      className={`p-6 rounded-2xl flex flex-col items-center justify-center gap-2 ${bg}`}
    >
      <span className={`text-3xl font-bold ${color}`}>{value}</span>
      <span className="text-sm font-medium text-gray-600 dark:text-gray-300">
        {label}
      </span>
    </div>
  );
}
