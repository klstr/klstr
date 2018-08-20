package command_jobs

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
)

type MySQLCommandJob struct {
	options CommandJobOptions
}

func (mcj MySQLCommandJob) BuildCommand(object *batchv1.Job) {
	object.Spec.Template.Spec.Containers[0].Image = "mysql"
	command := []string{
		"/bin/bash",
		"-c",
		"mysql",
		"--host=$MYSQLHOST",
		"--port=$MYSQLPORT",
		"--user=$MYSQLUSERNAME",
		"--password=$MYSQLPASSWORD",
		fmt.Sprintf("--execute='create database %s'", mcj.options.ToDBName),
		"&&",
		fmt.Sprintf("mysqldump %s", mcj.options.FromDBName),
		"|",
		fmt.Sprintf("mysql %s", mcj.options.ToDBName),
	}
	object.Spec.Template.Spec.Containers[0].Command = command
	secretKeyName := fmt.Sprintf("%s-%s-%s", mcj.options.DBIName, "mysql", mcj.options.FromDBName)
	env := []corev1.EnvVar{
		{
			Name: "MYSQLHOST",
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
			Name: "MYSQLPORT",
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
			Name: "MYSQLUSERNAME",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "username",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretKeyName,
					},
				},
			},
		},
		{
			Name: "MYSQLPASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: "password",
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretKeyName,
					},
				},
			},
		},
	}
	object.Spec.Template.Spec.Containers[0].Env = env
}

var _ CommandJob = MySQLCommandJob{}

func NewMySQLCommandJob(options CommandJobOptions) CommandJob {
	return &MySQLCommandJob{
		options: options,
	}
}
