package command_jobs

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type PGCommandJob struct {
	options CommandJobOptions
}

var _ CommandJob = PGCommandJob{}

func (pgcj PGCommandJob) getJobEnv() []corev1.EnvVar {
	secretKeyName := fmt.Sprintf("%s-%s-%s", pgcj.options.DBIName, "pg", pgcj.options.DBName)
	return []corev1.EnvVar{
		{
			Name: "PGHOST",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "host",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretKeyName,
					},
				},
			},
		},
		{
			Name: "PGPORT",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "port",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretKeyName,
					},
				},
			},
		},
		{
			Name: "PGPASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretKeyName,
					},
				},
			},
		},
		{
			Name: "PGUSERNAME",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "username",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretKeyName,
					},
				},
			},
		},
	}
}

func (pgcj PGCommandJob) getJobCommand(command []string) []string {
	return append([]string{
		"psql",
		"--host=$PGHOST",
		"--port=$PGPORT",
		"--username=$PGUSERNAME",
	}, command...)
}

func (pgcj PGCommandJob) BuildCloneCommand(object *batchv1.Job) {
	object.Spec.Template.Spec.Containers[0].Image = "postgres"
	cmd := []string{
		fmt.Sprintf(
			"--command='create database %s with template=%s'",
			pgcj.options.ToDBName,
			pgcj.options.DBName,
		),
	}
	object.Spec.Template.Spec.Containers[0].Command = pgcj.getJobCommand(cmd)
	object.Spec.Template.Spec.Containers[0].Env = pgcj.getJobEnv()
}

func (pgcj PGCommandJob) BuildCreateCommand(object *batchv1.Job) {
	object.Spec.Template.Spec.Containers[0].Image = "postgres"
	cmd := []string{
		fmt.Sprintf(
			"--command='create database %s'",
			pgcj.options.DBName,
		),
	}
	object.Spec.Template.Spec.Containers[0].Command = pgcj.getJobCommand(cmd)
	object.Spec.Template.Spec.Containers[0].Env = pgcj.getJobEnv()
}

func NewPGCommandJob(options CommandJobOptions) CommandJob {
	return &PGCommandJob{
		options: options,
	}
}
