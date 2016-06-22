package main

import (
	"github.com/FaXaq/gjp"
	"github.com/kataras/iris"
)

func TestPing(c *iris.Context) {
	c.JSON(iris.StatusOK, map[string]string{"ping": "pong"})
}

func CreateJob(c *iris.Context, jobPool *gjp.JobPool) (err error) {
	if len(c.URLParams()) > 0 {

		//creating custom job
		newJob, err := NewJob(
			gjp.GenerateJobUID(),
			c.URLParam("command"),
			c.URLParam("fromFile"),
			c.URLParam("toFile"))

		//if error while creating custom job, then print it in answer
		if err != nil {
			c.JSON(iris.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		} else {
			//if no error, queue the new job
			j := jobPool.QueueJob(newJob.Id, newJob.Name, newJob, 0)

			//get infos from the job
			c.JSON(iris.StatusOK, map[string]*gjp.Job{
				"job": j,
			})
		}
	} else {
		c.JSON(iris.StatusBadRequest, map[string]string{
			"error": "missing parameters",
		})
	}

	return
}

func GetMyJobProgress(c *iris.Context, jobPool *gjp.JobPool) (err error) {
	if len(c.URLParams()) > 0 &&
		c.URLParam("id") != "" {
		searchParam := c.URLParam("id")
		j, err := jobPool.GetJobFromJobId(searchParam)
		if err != nil {
			c.JSON(iris.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		} else {
			timing, err := j.GetProgress(j.GetJobId())
			if err != nil {
				c.JSON(iris.StatusInternalServerError, map[string]string{
					"error": err.Error(),
				})
			} else {
				//get infos from the job
				c.JSON(iris.StatusOK, map[string]float64{
					"percentage": timing,
				})
			}
		}
	} else {
		c.JSON(iris.StatusBadRequest, map[string]string{
			"error": "missing parameter",
		})
	}

	return
}

func SearchJob(c *iris.Context, jobPool *gjp.JobPool) (err error) {
	if len(c.URLParams()) > 0 {
		searchParam := c.URLParam("id")
		j, err := jobPool.GetJobFromJobId(searchParam)
		if err != nil {
			c.JSON(iris.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
		} else {
			c.JSON(iris.StatusOK, map[string]*gjp.Job{
				"job": j,
			})
		}
	} else {
		c.JSON(iris.StatusBadRequest, map[string]string{
			"error": "missing parameter",
		})
	}
	return
}

//NYI
func ListJobs(c *iris.Context, jobPool *gjp.JobPool) (err error) {
	c.JSON(iris.StatusOK, map[string]*gjp.JobPool{
		"jobs": jobPool,
	})

	return
}
