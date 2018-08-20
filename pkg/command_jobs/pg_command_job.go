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

func (pgcj PGCommandJob) BuildCommand(object *batchv1.Job) {
	object.Spec.Template.Spec.Containers[0].Image = "postgres"
	command := []string{
		"psql",
		"--host=$PGHOST",
		"--port=$PGPORT",
		"--username=$PGUSERNAME",
		fmt.Sprintf("--command='create database %s with template=%s'", pgcj.options.ToDBName, pgcj.options.FromDBName),
	}
	object.Spec.Template.Spec.Containers[0].Command = command
	secretKeyName := fmt.Sprintf("%s-%s-%s", pgcj.options.DBIName, "pg", pgcj.options.FromDBName)
	env := []corev1.EnvVar{
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
	object.Spec.Template.Spec.Containers[0].Env = env
}

var _ CommandJob = PGCommandJob{}

func NewPGCommandJob(options CommandJobOptions) CommandJob {
	return &PGCommandJob{
		options: options,
	}
}
