package models

import "fmt"

type CalculateImpactPointsParams struct {
	Contract Contract
	Project  Project
	Category JobCategory
}

func calculate(params CalculateImpactPointsParams) float64 {
	const RATIO float64 = 0.1
	experienceRatioMap := map[int]float64{
		0: -0.3,
		1: -0.1,
		2: 0,
		3: 0.1,
	}

	contract := params.Contract
	project := params.Project
	category := params.Category

	hourlyWage := *category.HourlyWageDollars
	totalHours := float64(contract.Commitment)                      //handle if null
	experienceRatio := experienceRatioMap[*project.ExperienceLevel] //handle if null

	totalPoints := hourlyWage * totalHours * (1 + experienceRatio)
	ratioPoints := totalPoints * RATIO

	if project.PaymentType != nil && *project.PaymentType == PaymentTypePaid {
		return ratioPoints
	}
	return totalPoints + ratioPoints
}

func CalculateImpactPoints(c *Contract) (float64, *CalculateImpactPointsParams, error) {
	contract, err := GetContract(c.ID)
	if err != nil {
		return 0, nil, err
	}

	//Getting Contract's project & validate
	project, err := GetProject(*contract.ProjectID)
	if err != nil {
		return 0, nil, err
	}

	if *project.PaymentScheme == PaymentSchemeFixed && contract.Status != ContractStatusCompleted {
		return 0, nil, fmt.Errorf("contract is not confirmed")
	}

	//Getting Project's category & validate
	if project.JobCategoryId == nil {
		return 0, nil, fmt.Errorf("there are no job category for project: %s", project.ID)
	}
	category, err := GetJobCategory(*project.JobCategoryId)
	if err != nil {
		return 0, nil, fmt.Errorf("there are no job category for project: %s", project.ID)
	}

	totalHours := float64(contract.Commitment)
	// if project.PaymentScheme == PaymentSchemeFixed {
	// 	totalHours = float64(contract.Commitment)
	// } else {
	// 	//TODO: Fix this if we currently use the submitted_work_id because it was all null
	// 	// Hourly basis jobs
	// 	// export const impactPointsCalculatedWorksIds = async (missionId) => {
	// 	// const { rows } = await app.db.query(sql`
	// 	// 	SELECT submitted_work_id FROM impact_points_history WHERE mission_id=${missionId}
	// 	// `)
	// 	// return rows.map((r) => r.submitted_work_id)
	// 	// }
	// 	// const calculatedWorks = await impactPointsCalculatedWorksIds(mission.id)
	// 	// const works = mission.submitted_works?.filter((w) => !calculatedWorks.includes(w.id))

	// 	// if (!works || works.length < 1) {
	// 	// 	logger.error(
	// 	// 	`Faild canculate impact score for ${mission.id},
	// 	// 	there are no submitted works for this mission that exists or not calculated already`
	// 	// 	)
	// 	// 	return false
	// 	// }
	// 	// workId = works[0].id
	// 	// totalHours = works[0].total_hours

	// 	totalHours = float64(contract.Commitment)
	// }

	if totalHours < 1 {
		return 0, nil, fmt.Errorf("total hours is under 1 hour")
	}

	calculateParams := CalculateImpactPointsParams{
		Contract: *contract,
		Project:  *project,
		Category: *category,
	}

	return calculate(calculateParams), &calculateParams, nil
}
